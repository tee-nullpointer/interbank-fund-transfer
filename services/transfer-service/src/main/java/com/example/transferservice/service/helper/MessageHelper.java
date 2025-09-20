package com.example.transferservice.service.helper;

import lombok.experimental.UtilityClass;
import org.springframework.util.StringUtils;

@UtilityClass
public class MessageHelper {

    public MessageType getMessageType(String processingCode) {
        if (!StringUtils.hasText(processingCode) || processingCode.length() != 6) {
            return null;
        }
        String messageType = processingCode.substring(0, 2);
        return switch (messageType) {
            case "43" -> MessageType.INQUIRY;
            case "91" -> MessageType.TRANSFER;
            default -> null;
        };
    }

    public InstrumentType[] getInstrumentType(String processingCode) {
        if (!StringUtils.hasText(processingCode) || processingCode.length() != 6) {
            return null;
        }
        String sourceType = processingCode.substring(2, 4);
        String destinationType = processingCode.substring(4, 6);
        InstrumentType sourceInstrumentType = switch (sourceType) {
            case "00" -> InstrumentType.CARD;
            case "20" -> InstrumentType.ACCOUNT;
            default -> null;
        };
        InstrumentType destinationInstrumentType = switch (destinationType) {
            case "00" -> InstrumentType.CARD;
            case "20" -> InstrumentType.ACCOUNT;
            default -> null;
        };
        return new InstrumentType[]{sourceInstrumentType, destinationInstrumentType};
    }
}
