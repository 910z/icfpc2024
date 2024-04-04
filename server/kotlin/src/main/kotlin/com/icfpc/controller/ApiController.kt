package com.icfpc.controller

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.ResponseBody

@Controller
@RequestMapping("/api")
class ApiController {
    @GetMapping("/problems")
    @ResponseBody
    fun problems(): List<String> {
        return listOf("abc","def","ghi")
    }
}