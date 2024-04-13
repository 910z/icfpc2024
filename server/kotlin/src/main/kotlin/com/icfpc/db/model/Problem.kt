package com.icfpc.db.model

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Id
import jakarta.persistence.Table

@Entity
@Table(name = "problem")
class Problem(
    @Id
    @Column
    val id: Int,
    @Column
    val contentId: Int,
    @Transient
    var bestSolution: Solution? = null
//    @Column
//    @JdbcTypeCode(SqlTypes.JSON)
//    val meta: Meta = Meta(),
//    @Column
//    var bestSolutionId: Int? = null
)

data class Meta(
    var instrs: Int? = null,
    var musicns: Int? = null,
//    val attends: Int? = null,
    var tastes: Int? = null,
    var pillars: Int? = null,
    var stageSize: String? = null
)

