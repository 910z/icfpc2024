package com.icfpc.db.repository

import com.icfpc.db.model.History
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query

interface HistoryRepository : JpaRepository<History, String> {
    @Query("SELECT history FROM History history ORDER BY history.createdAt desc LIMIT 100")
    fun latest(): List<History>
}
