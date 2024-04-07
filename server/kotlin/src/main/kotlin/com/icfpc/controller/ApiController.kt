package com.icfpc.controller

import com.fasterxml.jackson.databind.json.JsonMapper
import com.icfpc.db.model.Content
import com.icfpc.db.model.Problem
import com.icfpc.db.model.Solution
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import com.icfpc.db.repository.SolutionRepository
import com.icfpc.problem.model.Solve
import com.icfpc.utils.Json
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseBody
import java.io.File

@Controller
@RequestMapping("/api")
class ApiController(
    val problemRepository: ProblemRepository,
    val solutionRepository: SolutionRepository,
    val contentRepository: ContentRepository
) {
    @GetMapping("/init")
    fun init(): String {
        File("problems/old").listFiles()?.forEach {
            if (it.name.endsWith(".json")) {
                val id = it.name.substringBeforeLast(".json").toInt()
                val content = contentRepository.save(Content(content = it.readText().let { JsonMapper().readTree(it) }))
                problemRepository.save(Problem(id = id, contentId = content.id!!))
            }
        }

        return "redirect:/api/problems"
    }

    @GetMapping("/problems")
    @ResponseBody
    fun problems() = problemRepository.findAll().sortedBy { it.id }

    @GetMapping("/problem/{id}")
    @ResponseBody
    fun problem(@PathVariable id: Int) =
        contentRepository.getReferenceById(problemRepository.getReferenceById(id).contentId).content


    @GetMapping("/solutions")
    @ResponseBody
    fun solutions() = solutionRepository.findAll().sortedByDescending { it.id }

    @GetMapping("/solution/{id}")
    @ResponseBody
    fun solution(@PathVariable id: Int) =
        contentRepository.getReferenceById(solutionRepository.getReferenceById(id).contentId).content

    @PostMapping("/solution/{id}")
    @ResponseBody
    fun upload(@PathVariable id: Int, @RequestBody body: Solve): Solution {
        val content = contentRepository.save(Content(content = Json.toObject(body)))
        return solutionRepository.save(Solution(problemId = id, contentId = content.id!!))
    }

    @GetMapping("/test")
    @ResponseBody
    fun test() = solutionRepository.notCalculated("v3", 10)
}