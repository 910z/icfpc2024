package com.icfpc

import java.io.File

data class Point(val x: Int, val y: Int) {
    operator fun plus(p: Point) = Point(x + p.x, y + p.y)
    operator fun minus(p: Point) = Point(x - p.x, y - p.y)
    operator fun times(a: Int) = Point(a * x, a * y)
}

val steps = mapOf(
    Point(-1, -1) to 1,
    Point(0, -1) to 2,
    Point(1, -1) to 3,
    Point(-1, 0) to 4,
    Point(0, 0) to 5,
    Point(1, 0) to 6,
    Point(-1, 1) to 7,
    Point(0, 1) to 8,
    Point(1, 1) to 9,
)

val cacheMinD = mutableMapOf<Pair<Int, Int>, Int>()
val cache = mutableMapOf<Pair<Int, Int>, List<Int>>()

fun solve(file: String): String {
    val task = File(file).readText()
        .replace("\r", "")
        .substringBefore("\n\n")
        .split("\n")
        .map {
            val p = it.split(" ").map { it.toInt() }
            Point(p[0], p[1])
        }
        .distinct()
        .toMutableList()

    var pos = Point(0, 0)
    var speed = Point(0, 0)

    fun findMinD(v: Point, targets: List<Point>): Pair<Point,Int> {
        for (d in 0..10000) {
            val d2 = d * (d + 1) / 2
            val pp = pos + v * d
            for (s in targets) {
                val delta = s - pp
                if (Math.abs(delta.x) <= d2 && Math.abs(delta.y) <= d2) {
                    return s to d
                }
            }
        }
        error("d not found")
    }

    fun minD(v: Point, s: Point): Int {
        for (d in 0..10000) {
            val delta = s - v * d
            if (Math.abs(delta.x) <= d * (d + 1) / 2 && Math.abs(delta.y) <= d * (d + 1) / 2) {
                return d
            }
        }
        return 10000
    }

    fun buildPath(steps: Int, dd: Int): List<Int> {
            val list = mutableListOf<Int>()
            repeat(steps) {
                list.add(0)
            }
            var d = steps
            var dist = dd
            while (dist > 0 && d > 0) {
                if (dist >= d) {
                    dist -= d
                    list[steps - d]++
                }
                d--
            }
            return list
    }

    fun pathTo(target: Point): List<Point> {
        val ss = (target - pos)
        val d = minD(speed, ss)
        val dd = ss - speed * d
        val path1 = buildPath(d, Math.abs(dd.x)).let { list ->
            if (dd.x < 0) {
                list.map { -it }
            } else list
        }
        val path2 = buildPath(d, Math.abs(dd.y)).let { list ->
            if (dd.y < 0) {
                list.map { -it }
            } else list
        }
        val path = path1.mapIndexed { index, x ->
            val y = path2[index]
            Point(x, y)
        }
        return path
    }

    fun run(): String {
        val result = StringBuilder()
        val visited = mutableSetOf<Point>()
//        task.forEach { next ->
//            val start = listOf(pos, speed)
//            val path = pathTo(next)
//            path.forEach {
//                speed += it
//                pos += speed
//                result.append(steps[it])
//            }
//            if (pos != next) {
//                error("$pos != $next, $start")
//            }
//        }
        var stepCount = 0
        while (task.isNotEmpty()) {
//            val sorted = task.sortedBy { minD(speed, it - pos) }
            val next = task
                .let { t -> findMinD(speed, t).first }
                .let { it to pathTo(it) }
//                .filter { it !in visited }
//                .sortedBy { minD(speed, it - pos) }
//                    minD(speed, it - pos)
//                }
//                .take(1)
//                .map { it to pathTo(it) }
//                .minBy { it.second.size }
            visited += next.first
            val path = next.second
            path.forEach {
                speed += it
                pos += speed
                result.append(steps[it])
            }
            stepCount++
            if (stepCount % 10 == 0) {
                print(".")
            }
            if (stepCount % 100 == 0) {
                print(result.length)
            }
            task.remove(next.first)
        }
        println()
        return result.toString()
    }

    return run()
}

fun updOutput() {
    File("problems/spaceship/").let {
        it.listFiles()!!.filter { it.name.endsWith(".out") }
    }
        .map { it to it.readText().length }
        .sortedBy { it.second }
        .filter { it.second < 1_000_000 }
        .map { it.first }
        .joinToString("\n") {
        val task = it.absolutePath.substringBeforeLast(".").substringAfterLast("/").substringAfterLast("\\")
        "solve $task ${it.readText()}"
    }
        .let {
            File("problems/spaceship/output.txt").writeText(it)
        }
}

fun main() {
    try {
        updOutput()
        File("problems/spaceship/").let {
            it.listFiles()!!.filter { it.name.endsWith(".in") }
                .filter { !File(it.absolutePath.replace(".in", ".out")).exists() }
        }
            .shuffled()
            .sortedBy { it.readText().count { it == "\n"[0] } }
            .asSequence()
            .onEach { println(it.name) }
            .forEach {
                val task = it.absolutePath.substringBeforeLast(".in").substringAfterLast("/").substringAfterLast("\\")
                val solution = solve(it.absolutePath)
                File(it.absolutePath.replace(".in", ".out")).writeText(solution)
//            result.add("solve $task ${solve(it.absolutePath)}")
            }
    } finally {
        updOutput()
    }
}
