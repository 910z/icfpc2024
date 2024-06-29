package com.icfpc.problem.model

import com.fasterxml.jackson.annotation.JsonProperty
import com.fasterxml.jackson.databind.ObjectMapper
import java.math.BigDecimal
import java.text.DecimalFormat
import java.text.DecimalFormatSymbols
import java.util.*
import kotlin.math.sqrt

data class Solve(
    @JsonProperty("placements")
    val placements: List<Point>,
    @JsonProperty("volumes")
    val volumes: List<Double>? = (1..placements.size).map { 10.0 }
) {
    companion object {
        fun parse(json: String): Solve {
            return ObjectMapper().readValue(json, Solve::class.java)
        }
    }
}

data class Point(
    @JsonProperty("x")
    val x: Double,
    @JsonProperty("y")
    val y: Double
) {
    operator fun plus(p: Point) = Point(x + p.x, y + p.y)

    operator fun minus(p: Point) = Point(x - p.x, y - p.y)

    operator fun times(d: Double) = Point(x * d, y * d)

    operator fun div(d: Double) = Point(x / d, y / d)

    operator fun div(d: Int) = Point(x / d, y / d)

    infix fun scalar(p: Point) = p.x * x + p.y * y

    infix fun dist(p: Point) = sqrt((this - p).sqrSize())

    fun sqrSize() = x * x + y * y

    fun norm() = this / sqrt(sqrSize())

    override fun toString(): String = "${asString(x)} ${asString(y)}"
}

val df = DecimalFormat().apply {
    maximumFractionDigits = 3
    minimumFractionDigits = 0
    isGroupingUsed = false
    decimalFormatSymbols = DecimalFormatSymbols(Locale.ENGLISH)
}

fun asString(d: Double): String {
    return try {
         df.format(BigDecimal("$d"))
    } catch (e: Throwable) {
        "$d"
    }
}

//fun Solution.getContent(contentRepository: ContentRepository) =
//    contentRepository.getReferenceById(contentId).let { Json.parse<Solve>(it.content) }
