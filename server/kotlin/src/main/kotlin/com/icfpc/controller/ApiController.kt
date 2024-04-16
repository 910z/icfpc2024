package com.icfpc.controller

import com.fasterxml.jackson.databind.json.JsonMapper
import com.icfpc.db.model.Content
import com.icfpc.db.model.Problem
import com.icfpc.db.model.Solution
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import com.icfpc.db.repository.SolutionRepository
import com.icfpc.db.repository.bestSolutions
import com.icfpc.problem.model.Solve
import com.icfpc.utils.Json
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.*
import java.io.File

@RestController
@RequestMapping("/api")
class ApiController(
    val problemRepository: ProblemRepository,
    val solutionRepository: SolutionRepository,
    val contentRepository: ContentRepository
) {
    @GetMapping("/problems")
    fun problems(): List<Problem> {
        val problems = problemRepository.findAll().sortedBy { it.id }
        val best = solutionBest()
        problems.forEach {
            it.bestSolution = best[it.id]
        }
        return problems
    }

    fun solutionBest(): Map<Int, Solution> = solutionRepository.bestSolutions()

    @GetMapping("/problem/{id}")
    fun problem(@PathVariable id: Int) =
        contentRepository.getReferenceById(problemRepository.getReferenceById(id).contentId).content

    @GetMapping("/problem/tag")
    fun problemTag(tag: String) = problemRepository.findWithoutTag(tag)

    @GetMapping("/solutions")
    fun solutions(limit: Int = 50) = solutionRepository.findAll().sortedByDescending { it.id }.take(limit)

    @GetMapping("/solution/{id}")
    fun solution(@PathVariable id: Int) =
        contentRepository.getReferenceById(solutionRepository.getReferenceById(id).contentId).content

    @PostMapping("/solution/{id}")
    fun upload(@PathVariable id: Int, @RequestBody body: Solve, @RequestParam("tags") tags: List<String>): Solution {
        val content = contentRepository.save(Content(content = Json.toObject(body)))
        return solutionRepository.save(Solution(problemId = id, contentId = content.id!!, tags = tags))
    }

    @GetMapping("/test")
    fun test() = solutionRepository.notCalculated("v3", 10)
}