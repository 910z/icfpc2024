package com.icfpc.scheduler

import com.icfpc.db.model.Content
import com.icfpc.db.model.History
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.HistoryRepository
import com.icfpc.utils.QueryUtils
import org.springframework.scheduling.annotation.EnableScheduling
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Service

@Service
@EnableScheduling
class UpdateScoresService(
    val contentRepository: ContentRepository,
    val historyRepository: HistoryRepository,
    val queryUtils: QueryUtils
) {
    @Scheduled(fixedDelayString = "2000")
    fun update() {
        val history = queryUtils.getHistory()
        history.forEach {
            save(it)
        }
    }

    fun save(history: History) {
        if (!historyRepository.existsById(history.uuid)) {
            contentRepository.save(buildContent(history.request, queryUtils.getRequest(history.uuid)))
            contentRepository.save(buildContent(history.response, queryUtils.getResponse(history.uuid)))
            historyRepository.save(history)
        }
    }

    fun buildContent(id: String, content: String): Content {
        return Content(id, content, "")
    }
}