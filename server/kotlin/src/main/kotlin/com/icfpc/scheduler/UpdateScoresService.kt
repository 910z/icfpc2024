package com.icfpc.scheduler

import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import com.icfpc.db.repository.SolutionRepository
import com.icfpc.problem.CalcScoringService
import com.icfpc.problem.model.Solve
import com.icfpc.problem.model.Task
import com.icfpc.utils.Json
import org.springframework.scheduling.annotation.EnableScheduling
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Service
import java.util.*

@Service
@EnableScheduling
//@ConditionalOnProperty(name = ["service.calc"], matchIfMissing = false)
class UpdateScoresService(
    val calcMetric: CalcScoringService,
    val solutionRepository: SolutionRepository,
    val contentRepository: ContentRepository,
    val problemRepository: ProblemRepository
) {
    @Scheduled(fixedRateString = "1000")
    fun update() {
        val begin = Date()
        var first = true
        var res = true
        while (res && (Date().time - begin.time < 1000)) {
            res = calc()
            if (!res && first) {
                return
            }
            first = false
        }
//        println("score update end at ${Date().time - begin.time}ms")
    }

    fun calc(): Boolean {
        val version = CalcScoringService.currentVersion
        val solutions = solutionRepository.notCalculated(version, 5)
        if (solutions.isEmpty()) {
            return false
        }
        solutions.forEach { solution ->
            val begin = Date()
            val problem = problemRepository.getReferenceById(solution.problemId)
            val task = contentRepository.getReferenceById(problem.contentId).let { Json.parse<Task>(it.content) }
            val solve = contentRepository.getReferenceById(solution.contentId).let { Json.parse<Solve>(it.content) }
            val score = calcMetric.calc(task, solve).toBigInteger()
            solution.score = score
            solution.scores[version] = score
            solutionRepository.save(solution)
            println("calc ${solution.id}[${problem.id}] at ${Date().time - begin.time}ms")
        }
        return true
    }
}