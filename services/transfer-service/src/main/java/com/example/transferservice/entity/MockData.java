package com.example.transferservice.entity;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Entity
@Table(name = "mock_data", schema = "transfer")
public class MockData {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false)
    private Integer id;

    @Column(name = "instrument", length = 10)
    private String instrument;

    @Column(name = "number", length = 20)
    private String number;

    @Column(name = "inquiry_status")
    private String inquiryStatus;

    @Column(name = "name", length = 100)
    private String name;

    @Column(name = "transfer_status")
    private String transferStatus;

}