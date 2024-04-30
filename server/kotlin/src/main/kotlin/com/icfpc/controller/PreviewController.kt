package com.icfpc.controller

import com.github.nwillc.ksvg.RenderMode
import com.github.nwillc.ksvg.elements.SVG
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import com.icfpc.db.repository.SolutionRepository
import com.icfpc.problem.model.Point
import com.icfpc.problem.model.asString
import com.icfpc.problem.model.getContent
import org.springframework.cache.annotation.Cacheable
import org.springframework.http.CacheControl
import org.springframework.http.MediaType
import org.springframework.http.ResponseEntity
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.ResponseBody
import java.awt.Color
import java.awt.Graphics2D
import java.awt.RenderingHints
import java.awt.geom.Ellipse2D
import java.awt.image.BufferedImage
import java.io.File
import java.nio.file.Files
import java.util.concurrent.TimeUnit
import kotlin.math.max

@Controller
class PreviewController(
    val problemRepository: ProblemRepository,
    val solutionRepository: SolutionRepository,
    val contentRepository: ContentRepository
) {
    @GetMapping("/preview/{id}")
    @ResponseBody
    @Cacheable(value = ["preview"])//, key = "#id.#imgSize.#size")
    fun getImage(
        @PathVariable id: Int,
        imgSize: Int?,
        size: Int = 0
    ): ResponseEntity<ByteArray> {
        val solution = solutionRepository.getReferenceById(id)
        val problem = problemRepository.getReferenceById(solution.problemId)
        val task = problem.getContent(contentRepository)
        val solve = solution.getContent(contentRepository)

        val iSize = max(task.stage_width, task.stage_height) * 1.05
        val center = Point(
            task.stage_bottom_left[0] + task.stage_width / 2,
            task.stage_bottom_left[1] + task.stage_height / 2
        )

        val p0 = Point(-iSize, -iSize) * iSize / iSize / 2 + center
        val p1 = Point(iSize, iSize) * iSize / iSize / 2 + center

        val rect = Rectangle(p0, p1)


        val svgs = SVG.svg(true) {
            viewBox = "$rect"
            style {
                val m = task.musicians.max()
                val l = (0..m).map { convert(Color.getHSBColor(it.toFloat() / (m + 1), 0.5F, 0.5F)) }
                val bodyList = listOf(
                    "svg .p { fill:${convert(Color.LIGHT_GRAY)}; }",
                    "svg .a { fill:#FFFDD0; stroke:black; stroke-width:0.2; }"
                ) + l.mapIndexed { i, color -> "svg .m$i { fill:$color; }" }
                body = bodyList.joinToString("\n")
            }

            rect(
                Point(task.stage_bottom_left[0] + 5, task.stage_bottom_left[1] + 5),
                width = task.stage_width - 10,
                height = task.stage_height - 10,
                fill = Color.LIGHT_GRAY
            )

            solve.placements.forEachIndexed { index, it ->
                circle(rect, it, 5.0, "m${task.musicians[index]}")
            }

            task.pillars.forEach { pillar ->
                circle(rect, Point(pillar.center[0], pillar.center[1]), pillar.radius, "p")
            }

            task.attendees.forEach { attendee ->
                circle(rect, Point(attendee.x, attendee.y), 5.0, "a")
            }
        }

//        val svg = ImageDrawSVG(imgSize ?: 1000, center, iSize) {
//            svg.style {
//                val m = task.musicians.max()
//                val l = (0..m).map { convert(Color.getHSBColor(it.toFloat() / (m + 1), 0.5F, 0.5F)) }
//                val bodyList = listOf(
//                    "svg .p { fill:${convert(Color.LIGHT_GRAY)}; }",
//                    "svg .a { fill:#FFFDD0; stroke:black; stroke-width:0.2; }"
//                ) + l.mapIndexed { i, color -> "svg .m$i { fill:$color; }" }
//                body = bodyList.joinToString("\n")
//            }
//
//            fillRect(
//                Point(task.stage_bottom_left[0] + 5, task.stage_bottom_left[1] + 5),
//                task.stage_width - 10,
//                task.stage_height - 10,
//                Color.LIGHT_GRAY
//            )
//
//            solve.placements.forEachIndexed { index, it ->
//                fillCircle(it, 5.0, "m${task.musicians[index]}")
//            }
//
//            task.pillars.forEach { pillar ->
//                fillCircle(Point(pillar.center[0], pillar.center[1]), pillar.radius, "p")
//            }
//
//            task.attendees.forEach { attendee ->
//                fillCircle(Point(attendee.x, attendee.y), 5.0, "a")
//            }
//        }

        val res = StringBuilder().also {
            svgs.render(it, RenderMode.FILE)
        }.toString().toByteArray()

        return ResponseEntity.ok()
            .contentType(MediaType("image", "svg+xml"))
            .cacheControl(CacheControl.maxAge(1, TimeUnit.MINUTES))
            .body(res)
    }

    fun cache(key: String, lambda: () -> ByteArray): ByteArray {
        val file = File("$DIR/$key")
        if (!file.exists()) {
            file.parentFile.mkdirs()
            Files.write(file.toPath(), lambda())
        }
        return file.readBytes()
    }

    companion object {
        val DIR = "temp"
    }
}

data class ImageDraw(val size: Int, val center: Point, val scale: Double, val image: BufferedImage) {
    val g2d = image.graphics as Graphics2D

    var color: Color
        get() {
            TODO()
        }
        set(value) {
            g2d.color = value
        }

    init {
        color = Color(255, 0, 255, 0)
        g2d.fillRect(0, 0, size, size)
        g2d.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON)
        g2d.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON)
    }

    companion object {
        operator fun invoke(size: Int, center: Point, scale: Double, draw: ImageDraw.() -> Unit): BufferedImage {
            val image = BufferedImage(size, size, BufferedImage.TYPE_INT_ARGB)
            val drawI = ImageDraw(size, center, scale, image)
            drawI.draw()
            return drawI.image
        }
    }

    infix fun Point.pMul(a: Point) = Point(this.x * a.x, this.y * a.y)
    infix fun Point.pDiv(a: Point) = Point(this.x / a.x, this.y / a.y)

    fun convert(p: Point) = ((p - center) * size.toDouble() / scale) + Point(size.toDouble() / 2, size.toDouble() / 2)
    fun convert(d: Double) = d * size.toDouble() / scale

    fun fillCircle(p: Point, r: Double) {
        val a = convert(p)
        val rd = convert(r)

        val shape = Ellipse2D.Double(a.x - rd, a.y - rd, rd * 2, rd * 2)
        g2d.fill(shape)
//        image.setRGB(a.x.toInt(), a.y.toInt(), g2d.color.rgb)
    }

    fun fillRect(from: Point, width: Double, height: Double) {
        val a = convert(from)
        g2d.fillRect(a.x.toInt(), a.y.toInt(), convert(width).toInt(), convert(height).toInt())
    }

    fun drawRect(from: Point, width: Double, height: Double) {
        val a = convert(from)
        g2d.drawRect(a.x.toInt(), a.y.toInt(), convert(width).toInt(), convert(height).toInt())
    }
}

data class ImageDrawSVG(val size: Int, val center: Point, val scale: Double, val svg: SVG) {
//    var color: Color
//        get() {
//            TODO()
//        }
//        set(value) {
//            g2d.color = value
//        }

    init {
//        color = Color(255, 0, 255, 0)
//        g2d.fillRect(0, 0, size, size)
//        g2d.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON)
//        g2d.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON)
    }

    companion object {
        operator fun invoke(size: Int, center: Point, scale: Double, draw: ImageDrawSVG.() -> Unit): SVG {
            return SVG.svg(true) {
                val p0 = Point(-size.toDouble(), -size.toDouble()) * scale / size / 2 + center
                val p1 = Point(size.toDouble(), size.toDouble()) * scale / size / 2 + center
                viewBox = "$p0 ${p1 - p0}"
//                val iSize = max(task.stage_width, task.stage_height) * 1.05
//                val center = Point(
//                    task.stage_bottom_left[0] + task.stage_width / 2,
//                    task.stage_bottom_left[1] + task.stage_height / 2
//                )
//                fun convert(p: Point) = ((p - center) * size.toDouble() / scale) + Point(size.toDouble() / 2, size.toDouble() / 2)
//                fun convert(d: Double) = d * size.toDouble() / scale
//                x * scale  = d * size.toDouble() / scale
                val svg = this
                val drawI = ImageDrawSVG(size, center, scale, svg)
                drawI.draw()
            }
        }
    }

    infix fun Point.pMul(a: Point) = Point(this.x * a.x, this.y * a.y)
    infix fun Point.pDiv(a: Point) = Point(this.x / a.x, this.y / a.y)

    fun convert(p: Point) = ((p - center) * size.toDouble() / scale) + Point(size.toDouble() / 2, size.toDouble() / 2)
    fun convert(d: Double) = d * size.toDouble() / scale


    fun fillCircle(p: Point, r: Double, cssClass: String) {
        val a = convert(p)
        val rd = convert(r)

//        val shape = Ellipse2D.Double(a.x - rd, a.y - rd, rd * 2, rd * 2)
        if (a.x <= size + rd && a.y <= size + rd && a.x >= -rd && a.y >= -rd) {
            svg.circle {
//                cx = "${a.x}"
//                cy = "${a.y}"
                cx = asString(p.x)
                cy = asString(p.y)
                this.r = asString(r)
                this.cssClass = cssClass
            }
        }
//        g2d.fill(shape)
//        image.setRGB(a.x.toInt(), a.y.toInt(), g2d.color.rgb)
    }

    fun fillRect(from: Point, width: Double, height: Double, color: Color) {
        val a = convert(from)
        svg.rect {
            x = asString(from.x)
            y = asString(from.y)
            this.width = asString(width)
            this.height = asString(height)
            fill = convert(color)
        }
//        g2d.fillRect(a.x.toInt(), a.y.toInt(), convert(width).toInt(), convert(height).toInt())
    }

//    fun SVG.drawRect(from: Point, width: Double, height: Double) {
//        val a = convert(from)
////        g2d.drawRect(a.x.toInt(), a.y.toInt(), convert(width).toInt(), convert(height).toInt())
//    }
}

fun convert(c: Color): String {
    val buf = Integer.toHexString(c.rgb)
    return "#" + buf.substring(buf.length - 6)
}

data class Rectangle(val p0: Point, val p1: Point) {
    override fun toString() = "$p0 ${p1 - p0}"
}

fun SVG.rect(from: Point, width: Double, height: Double, fill: Color? = null) {
    rect {
        x = asString(from.x)
        y = asString(from.y)
        this.width = asString(width)
        this.height = asString(height)
        fill?.let {
            this.fill = convert(it)
        }
    }
}

fun SVG.circle(viewBox: Rectangle, p: Point, r: Double, cssClass: String) {
    if (viewBox.p0.x - r <= p.x && p.x <= viewBox.p1.x + r && viewBox.p0.y - r <= p.y && p.y <= viewBox.p1.y + r) {
        circle {
            cx = asString(p.x)
            cy = asString(p.y)
            this.r = asString(r)
            this.cssClass = cssClass
        }
    }
}
