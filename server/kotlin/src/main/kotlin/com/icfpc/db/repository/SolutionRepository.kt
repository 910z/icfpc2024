package com.icfpc.db.repository

import com.icfpc.db.model.Solution
import org.springframework.data.jpa.repository.JpaRepository

interface SolutionRepository : JpaRepository<Solution, Int>