package com.example.transferservice.util;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.experimental.UtilityClass;
import lombok.extern.slf4j.Slf4j;

import java.util.Optional;

@UtilityClass
@Slf4j
public class ObjectMapperUtil {
    private final ObjectMapper mapper = new ObjectMapper();

    public Optional<String> toJson(Object obj) {
        try {
            return Optional.of(mapper.writeValueAsString(obj));
        } catch (Exception e) {
            log.error("Failed to convert object to JSON", e);
        }
        return Optional.empty();
    }

    public <T> Optional<T> fromJson(String json, Class<T> clazz) {
        try {
            return Optional.of(mapper.readValue(json, clazz));
        } catch (Exception e) {
            log.error("Failed to convert JSON to object", e);
        }
        return Optional.empty();
    }
}
