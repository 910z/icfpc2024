package com.icfpc.db.model

import com.fasterxml.jackson.databind.JsonNode
import jakarta.persistence.*
import org.hibernate.annotations.JdbcTypeCode
import org.hibernate.type.SqlTypes

@Entity
@Table(name = "content")
class Content(
    @Id
    @GeneratedValue
    @Column
    val id: Int? = null,
    @Column
    @JdbcTypeCode(SqlTypes.JSON)
    val content: JsonNode
)
