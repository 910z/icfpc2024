package com.icfpc.controller

import com.icfpc.db.repository.ContentRepository
import com.icfpc.utils.QueryUtils
import org.springframework.stereotype.Controller

@Controller
class MainController(
    val contentRepository: ContentRepository,
    val queryUtils: QueryUtils
) {
//    @GetMapping("/history")
//    @ResponseBody
//    fun history() = queryUtils.history()

//    @GetMapping("/")
//    @ResponseBody
//    fun index() = this::class.java
//        .getResourceAsStream("/static/index.html")!!
//        .readAllBytes()
//
//    @GetMapping("/init")
//    fun init(): String {
//        (File("problems/old").listFiles() ?: emptyArray())
//            .mapNotNull {
//                if (it.name.endsWith(".json")) {
//                    val id = it.name.substringBeforeLast(".json").toInt()
//                    val content =
//                        contentRepository.save(Content(content = it.readText().let { JsonMapper().readTree(it) }))
//                    Problem(id = id, contentId = content.id!!)
//                } else {
//                    null
//                }
//            }
//            .sortedBy { it.id }
//            .forEach { problemRepository.save(it) }
//
//        return "redirect:/"
//    }
}
