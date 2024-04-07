package com.icfpc.db.repository

import com.icfpc.db.model.Problem
import org.springframework.data.jpa.repository.JpaRepository

interface ProblemRepository : JpaRepository<Problem, Int>
