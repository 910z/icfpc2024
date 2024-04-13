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

@Controller
@RequestMapping("/api")
class ApiController(
    val problemRepository: ProblemRepository,
    val solutionRepository: SolutionRepository,
    val contentRepository: ContentRepository
) {
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

        return "redirect:/api/problems"
    }

    @GetMapping("/problems")
    @ResponseBody
    fun problems(): List<Problem> {
        val problems = problemRepository.findAll().sortedBy { it.id }
        val best = solutionBest()
        problems.forEach {
            it.bestSolution = best[it.id]
        }
        return problems
    }

    //    @GetMapping("/solution/best")
    fun solutionBest(): Map<Int, Solution> = solutionRepository.bestSolutions()

    @GetMapping("/problem/{id}")
    @ResponseBody
    fun problem(@PathVariable id: Int) =
        contentRepository.getReferenceById(problemRepository.getReferenceById(id).contentId).content

    @GetMapping("/problem/tag")
    @ResponseBody
    fun problemTag(tag: String) = problemRepository.findWithoutTag(tag)

    @GetMapping("/solutions")
    @ResponseBody
    fun solutions(limit: Int = 50) = solutionRepository.findAll().sortedByDescending { it.id }.take(limit)

    @GetMapping("/solution/{id}")
    @ResponseBody
    fun solution(@PathVariable id: Int) =
        contentRepository.getReferenceById(solutionRepository.getReferenceById(id).contentId).content

    @PostMapping("/solution/{id}")
    @ResponseBody
    fun upload(@PathVariable id: Int, @RequestBody body: Solve, @RequestParam("tags") tags: List<String>): Solution {
        val content = contentRepository.save(Content(content = Json.toObject(body)))
        return solutionRepository.save(Solution(problemId = id, contentId = content.id!!, tags = tags))
    }

    @GetMapping("/test")
    @ResponseBody
    fun test() = solutionRepository.notCalculated("v3", 10)
}