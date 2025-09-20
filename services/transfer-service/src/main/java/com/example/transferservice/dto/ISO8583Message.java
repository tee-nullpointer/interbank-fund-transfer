package com.example.transferservice.dto;

import lombok.Data;

import java.util.Map;

@Data
public class ISO8583Message {
    private String mti;
    private Map<Integer, String> fields;


    public String describe() {
        StringBuilder sb = new StringBuilder();
        sb.append("MTI : ").append(this.mti);
        for (Map.Entry<Integer, String> entry : fields.entrySet()) {
            sb.append("\n");
            sb.append("F").append(entry.getKey()).append(" : ").append(entry.getValue());
        }
        return sb.toString();
    }
}
