package com.icfpc

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.boot.web.servlet.ServletComponentScan

@ServletComponentScan
@SpringBootApplication
class IcfpcApplication

fun main(args: Array<String>) {
	runApplication<IcfpcApplication>(*args)
}
