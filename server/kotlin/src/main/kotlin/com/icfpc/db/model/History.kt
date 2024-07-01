package com.icfpc.db.model

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.Id
import jakarta.persistence.Table


@Entity
@Table(name = "history")
class History(
    @Id
    @Column
    val uuid: String,
    @Column
    val createdAt: String,
    @Column
    val request: String,
    @Column
    val response: String
)
