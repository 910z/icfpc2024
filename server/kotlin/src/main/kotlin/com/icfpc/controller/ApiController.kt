package com.icfpc.controller

import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.HistoryRepository
import com.icfpc.utils.Json
import com.icfpc.utils.QueryUtils
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("api")
class ApiController(
    val historyRepository: HistoryRepository,
    val contentRepository: ContentRepository
) {
    @GetMapping("history", produces= ["application/json"])
    fun history(): Any {
        val history = historyRepository.latest()
        val content = history.flatMap {
            listOf(it.request, it.response)
        }
            .distinct()
            .associateWith { contentRepository.getReferenceById(it) }
        return mapOf(
            "history" to history,
            "content" to content
        )
    }
}