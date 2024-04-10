package com.icfpc.db.repository

import com.icfpc.db.model.Problem
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query

interface ProblemRepository : JpaRepository<Problem, Int> {
    @Query(
        "select p.id from problem p where p.id not in (" +
                "select s.problem_id from solution s inner join solution_tags st on s.id = st.solution_id where st.tags = ?1" +
                ") order by p.id asc", nativeQuery = true
    )
    fun findWithoutTag(tag: String): List<Int>
}
