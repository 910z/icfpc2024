package com.icfpc.db.model

import jakarta.persistence.*
import org.hibernate.annotations.JdbcTypeCode
import org.hibernate.type.SqlTypes
import java.math.BigInteger

typealias Score = MutableMap<String, BigInteger>

@Entity
@Table(name = "solution")
class Solution(
    @Id
    @GeneratedValue
    @Column
    val id: Int? = null,
    @Column
    val problemId: Int,
    @Column
    val contentId: Int,
//    @Column
//    val tag: String,
    @Column
    @ElementCollection
    var tags: List<String> = listOf(),
    @Column
    var score: BigInteger? = null,
    @Column
    @JdbcTypeCode(SqlTypes.JSON)
    val scores: Score = mutableMapOf()
)
