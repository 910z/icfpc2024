package com.icfpc.db.repository

import com.icfpc.db.model.Content
import org.springframework.data.jpa.repository.JpaRepository

interface ContentRepository : JpaRepository<Content, String>