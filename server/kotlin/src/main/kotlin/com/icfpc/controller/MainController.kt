package com.icfpc.controller

import com.fasterxml.jackson.databind.json.JsonMapper
import com.icfpc.db.model.Content
import com.icfpc.db.model.Problem
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.ResponseBody
import java.io.File

@Controller
class MainController(
    val problemRepository: ProblemRepository,
    val contentRepository: ContentRepository
) {
    @GetMapping("/")
    @ResponseBody
    fun index() = this::class.java
        .getResourceAsStream("/static/index.html")!!
        .readAllBytes()

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
