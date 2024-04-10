package com.icfpc.db.repository

import com.icfpc.db.model.Problem
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query

interface ProblemRepository : JpaRepository<Problem, Int> {
    @Query(
        "select p.id from problem p where p.id not in (" +
                "select s.problem_id from solution s where ?1=any(s.tags)" +
                ") order by p.id asc", nativeQuery = true
    )
    fun findWithoutTag(tag: String): List<Int>
}
