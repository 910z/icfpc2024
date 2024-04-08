package com.icfpc.controller

import jakarta.servlet.*
import jakarta.servlet.annotation.WebFilter
import jakarta.servlet.http.HttpServletResponse

@WebFilter("/*")
class AddResponseHeaderFilter : Filter {
    override fun doFilter(
        request: ServletRequest?,
        response: ServletResponse,
        chain: FilterChain
    ) {
        val httpServletResponse = response as HttpServletResponse
        httpServletResponse.setHeader("Access-Control-Allow-Origin", "*")
        chain.doFilter(request, response)
    }
}