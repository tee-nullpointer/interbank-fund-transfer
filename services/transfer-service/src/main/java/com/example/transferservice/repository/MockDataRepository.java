package com.example.transferservice.repository;

import com.example.transferservice.entity.MockData;
import org.springframework.data.jpa.repository.JpaRepository;

public interface MockDataRepository extends JpaRepository<MockData, Integer> {
}
