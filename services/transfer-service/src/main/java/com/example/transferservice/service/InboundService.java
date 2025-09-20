package com.example.transferservice.service;

import com.example.transferservice.dto.ISO8583Message;

public interface InboundService {
    void process(ISO8583Message message, String serviceId, String traceId);
}
