package com.example.transferservice.consumer;

import com.example.transferservice.dto.ISO8583Message;
import com.example.transferservice.service.InboundService;
import com.example.transferservice.util.ObjectMapperUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.util.Optional;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

@Component
@RequiredArgsConstructor
@Slf4j
public class InboundRequestConsumer {
    private final ExecutorService executor = Executors.newVirtualThreadPerTaskExecutor();
    private final InboundService inboundService;

    @KafkaListener(topics = "${application.kafka.request-topic}", groupId = "${application.kafka.consumer-group-id}")
    public void listen(String jsonMessage,
                       @Header(name = "service_id") String serviceId,
                       @Header(name = "trace_id") String traceId) {
        log.info("Received message: {}", jsonMessage);
        Optional<ISO8583Message> optional = ObjectMapperUtil.fromJson(jsonMessage, ISO8583Message.class);
        if (optional.isEmpty()) {
            return;
        }
        ISO8583Message message = optional.get();
        if (!StringUtils.hasText(serviceId)) {
            log.warn("Message missing service_id header, skipping processing");
            return;
        }
        executor.submit(() -> inboundService.process(message, serviceId, traceId));
    }
}
