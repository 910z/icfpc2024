//package com.icfpc.db.repository
//
//import com.icfpc.db.model.History
//import org.springframework.data.jpa.repository.JpaRepository
//
//interface SolutionRepository : JpaRepository<History, Int> {
////    @Query("select * from solution where scores -> ?1 is null order by id desc limit ?2", nativeQuery = true)
////    fun notCalculated(version: String, limit: Int): List<Solution>
////
////    @Query("select distinct on (solution.problem_id) solution.* " +
////            "from (" +
////            "select problem_id, max(score) as score " +
////            "from solution " +
////            "group by problem_id" +
////            ") maxs " +
////            "inner join solution on solution.problem_id = maxs.problem_id and solution.score = maxs.score", nativeQuery = true)
////    fun bestSolutionList(): List<Solution>
////
////    fun findFirstByProblemIdAndScoreIsNotNullOrderByScoreDescIdAsc(): List<Solution>
//}
//
////fun SolutionRepository.bestSolutions() = bestSolutionList()
////    .groupBy { it.problemId }
////    .mapValues { it.value.first() }
//////findAll()
//////    .groupBy { it.problemId }
//////    .mapValues { list -> list.value.maxBy { it.score ?: BigInteger.ZERO } }
