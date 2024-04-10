package com.icfpc.db.model

import com.icfpc.db.repository.ContentRepository
import com.icfpc.problem.model.Task
import com.icfpc.utils.Json
import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Id
import jakarta.persistence.Table
import java.math.BigInteger

@Entity
@Table(name = "problem")
class Problem(
    @Id
    @Column
    val id: Int,
    @Column
    val contentId: Int
)

