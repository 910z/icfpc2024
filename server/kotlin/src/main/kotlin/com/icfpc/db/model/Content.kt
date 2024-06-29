package com.icfpc.db.model

import jakarta.persistence.*
import org.hibernate.annotations.JdbcTypeCode
import java.sql.Types

@Entity
@Table(name = "content")
class Content(
    @Id
    @Column
    val id: String,
    @Column(columnDefinition = "TEXT")
    val content: String,
    @Column(columnDefinition = "TEXT")
    val parsed: String
)
