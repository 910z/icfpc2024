package com.icfpc.db.repository

import com.icfpc.db.model.Solution
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query
import java.lang.annotation.Native

interface SolutionRepository : JpaRepository<Solution, Int> {
    @Query("select * from solution where scores -> ?1 is null order by id desc limit ?2", nativeQuery = true)
    fun notCalculated(version: String, limit: Int): List<Solution>
}