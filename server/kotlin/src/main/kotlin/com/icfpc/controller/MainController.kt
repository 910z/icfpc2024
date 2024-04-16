package com.icfpc.controller

import com.fasterxml.jackson.databind.json.JsonMapper
import com.icfpc.db.model.Content
import com.icfpc.db.model.Problem
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import com.icfpc.db.repository.SolutionRepository
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import java.io.File

@Controller
class MainController(
    val problemRepository: ProblemRepository,
    val contentRepository: ContentRepository
) {
    @GetMapping("/")
    fun index(): String {
        return "redirect:index.html"
    }

    @GetMapping("/init")
    fun init(): String {
        (File("problems/old").listFiles() ?: emptyArray())
            .mapNotNull {
                if (it.name.endsWith(".json")) {
                    val id = it.name.substringBeforeLast(".json").toInt()
                    val content =
                        contentRepository.save(Content(content = it.readText().let { JsonMapper().readTree(it) }))
                    Problem(id = id, contentId = content.id!!)
                } else {
                    null
                }
            }
            .sortedBy { it.id }
            .forEach { problemRepository.save(it) }

        return "redirect:/"
    }
}