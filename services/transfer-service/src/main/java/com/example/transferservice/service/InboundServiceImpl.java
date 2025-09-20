package com.example.transferservice.service;

import com.example.transferservice.dto.ISO8583Message;
import com.example.transferservice.entity.MockData;
import com.example.transferservice.repository.MockDataRepository;
import com.example.transferservice.service.helper.InstrumentType;
import com.example.transferservice.service.helper.MessageHelper;
import com.example.transferservice.service.helper.MessageType;
import com.example.transferservice.util.ObjectMapperUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.Message;
import org.springframework.messaging.support.MessageBuilder;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class InboundServiceImpl implements InboundService {
    private final MockDataRepository mockDataRepository;
    private final KafkaTemplate<String, String> kafkaTemplate;
    private List<MockData> mockDataList = new ArrayList<>();
    @Value("${application.kafka.response-topic}")
    private String responseTopic;

    @Override
    public void process(ISO8583Message message, String serviceId, String traceId) {
        setupTraceLog(traceId);
        try {
            _process(message, serviceId);
        } catch (Exception e) {
            log.error("Fail to process inbound message", e);
        } finally {
            MDC.clear();
        }
    }

    @Scheduled(fixedDelay = 60000)
    private void getMockData() {
        mockDataList = mockDataRepository.findAll();
    }

    private void _process(ISO8583Message message, String serviceId) {
        log.info("Start processing message : {}", message.describe());
        setupResponse(message);
        sendResponse(message, serviceId);
    }

    private void setupResponse(ISO8583Message message) {
        String processingCode = message.getFields().get(3);
        String toAccount = message.getFields().get(103);
        if (!StringUtils.hasText(toAccount)) {
            log.warn("To account is empty");
            message.getFields().put(39, "30");
            return;
        }
        MessageType messageType = MessageHelper.getMessageType(processingCode);
        InstrumentType[] instrumentTypes = MessageHelper.getInstrumentType(processingCode);
        if (messageType == null || instrumentTypes == null) {
            log.warn("Invalid processing code : {}", processingCode);
            message.getFields().put(39, "30");
            return;
        }
        MockData mockData = mockDataList
                .stream()
                .filter(m -> toAccount.equals(m.getNumber()) &&
                        instrumentTypes[1] == InstrumentType.valueOf(m.getInstrument()))
                .findFirst()
                .orElse(null);
        if (mockData == null) {
            log.warn("To account not found : {}", toAccount);
            message.getFields().put(39, "14");
            return;
        }
        switch (messageType) {
            case INQUIRY: {
                message.getFields().put(39, mockData.getInquiryStatus());
                if ("00".equals(mockData.getInquiryStatus())) {
                    message.getFields().put(120, mockData.getName());
                }
                break;
            }
            case TRANSFER: {
                message.getFields().put(39, mockData.getTransferStatus());
                break;
            }
        }
    }

    private void sendResponse(ISO8583Message message, String serviceId) {
        Optional<String> optional = ObjectMapperUtil.toJson(message);
        if (optional.isEmpty()) {
            return;
        }
        Message<String> kafkaMessage = MessageBuilder
                .withPayload(optional.get())
                .setHeader(KafkaHeaders.TOPIC, responseTopic)
                .setHeader(KafkaHeaders.KEY, serviceId)
                .setHeader("service_id", serviceId)
                .setHeader("trace_id", MDC.get("traceId"))
                .build();
        kafkaTemplate.send(kafkaMessage);
        log.info("Sent response message : {} to topic : {}", message.describe(), responseTopic);
    }

    private void setupTraceLog(String traceId) {
        MDC.put("traceId", traceId);
        MDC.put("spanId", UUID.randomUUID().toString());
    }


}
