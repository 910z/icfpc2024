package com.icfpc.controller

import com.github.nwillc.ksvg.RenderMode
import com.github.nwillc.ksvg.elements.SVG
import com.icfpc.db.repository.ContentRepository
import com.icfpc.db.repository.ProblemRepository
import com.icfpc.db.repository.SolutionRepository
import com.icfpc.problem.model.Point
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
import java.io.ByteArrayOutputStream
import java.io.File
import java.io.FileWriter
import java.nio.file.Files
import java.util.concurrent.TimeUnit
import javax.imageio.ImageIO
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

        val res = if (size == 1) {
            val iSize = max(task.attendees.map { it.x }.max() + 5.0, task.attendees.map { it.y }.max() + 5.0)
            val center = Point(
                (task.attendees.map { it.x }.max() + 5.0) / 2,
                (task.attendees.map { it.y }.max() + 5.0) / 2
            )
            val image = ImageDraw(imgSize ?: 1000, center, iSize) {
                color = Color.LIGHT_GRAY
                fillRect(
                    Point(task.stage_bottom_left[0], task.stage_bottom_left[1]),
                    task.stage_width,
                    task.stage_height,
                )
                color = Color.BLACK
                drawRect(
                    Point(task.stage_bottom_left[0], task.stage_bottom_left[1]),
                    task.stage_width,
                    task.stage_height,
                )

                solve.placements.forEachIndexed { index, it ->
                    color = Color.getHSBColor(
                        task.musicians[index].toFloat() / task.musicians.max(),
                        0.5F,
                        0.5F
                    )
                    fillCircle(it, 5.0)
                }

                color = Color.LIGHT_GRAY

                task.pillars.forEach { pillar ->
                    fillCircle(Point(pillar.center[0], pillar.center[1]), pillar.radius)
                }

                color = Color(0xFF, 0xFD, 0xD0)

                task.attendees.forEach { attendee ->
                    fillCircle(Point(attendee.x, attendee.y), 5.0)
                }
            }

            val baos = ByteArrayOutputStream()
            ImageIO.write(image, "PNG", baos)
            baos.toByteArray()
        } else {
            val iSize = max(task.stage_width, task.stage_height) * 1.05
            val center = Point(
                task.stage_bottom_left[0] + task.stage_width / 2,
                task.stage_bottom_left[1] + task.stage_height / 2
            )

//                val s = SVG.svg(true) {
//                    val iSize = max(task.stage_width, task.stage_height) * 1.05
//                    val center = Point(
//                        task.stage_bottom_left[0] + task.stage_width / 2,
//                        task.stage_bottom_left[1] + task.stage_height / 2
//                    )
//                    this.viewBox = "${center - Point(iSize, iSize)/2} ${Point(iSize, iSize)}"
//
//
//                }

            val svg = ImageDrawSVG(imgSize ?: 1000, center, iSize) {
                svg.style {
                    val m = task.musicians.max()
                    val l = (0..m).map { convert(Color.getHSBColor(it.toFloat() / (m + 1), 0.5F, 0.5F)) }
                    body = """
svg .p { fill:${convert(Color.LIGHT_GRAY)}; }
svg .a { fill:#FFFDD0; stroke:black; stroke-width:0.2; }
${l.mapIndexed { i, color -> "svg .m$i { fill:$color; }" }.joinToString("\n")}
                        """.trimIndent()
                }

//                svg .a { fill:${convert(Color(0, 2, 0x2F))}; }
//                body[data-bs-theme=dark] svg .a { fill:${convert(Color(0xFF, 0xFD, 0xD0))};  }
//                @media (prefers-color-scheme: dark) { svg .a { fill:${convert(Color(0xFF, 0xFD, 0xD0))}; } }

//                @media(prefers-color-scheme: dark) svg .a { fill:'${convert(Color(0xFF, 0xFD, 0xD0))}'; }

                fillRect(
                    Point(task.stage_bottom_left[0] + 5, task.stage_bottom_left[1] + 5),
                    task.stage_width - 10,
                    task.stage_height - 10,
                    Color.LIGHT_GRAY
                )

                solve.placements.forEachIndexed { index, it ->
//                        val color = Color.getHSBColor(
//                            task.musicians[index].toFloat() / task.musicians.max(),
//                            0.5F,
//                            0.5F
//                        )в
                    fillCircle(it, 5.0, "m${task.musicians[index]}")
                }

                task.pillars.forEach { pillar ->
                    fillCircle(Point(pillar.center[0], pillar.center[1]), pillar.radius, "p")
                }

                task.attendees.forEach { attendee ->
                    fillCircle(Point(attendee.x, attendee.y), 5.0, "a")
                }
            }

            synchronized(DIR) {
                FileWriter("temp.svg").use {
                    svg.render(it, RenderMode.FILE)
                }
                File("temp.svg").readBytes()
            }
        }
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
                cx = "${p.x}"
                cy = "${p.y}"
                this.r = "$r"
                this.cssClass = cssClass
            }
        }
//        g2d.fill(shape)
//        image.setRGB(a.x.toInt(), a.y.toInt(), g2d.color.rgb)
    }

    fun fillRect(from: Point, width: Double, height: Double, color: Color) {
        val a = convert(from)
        svg.rect {
            x = "${from.x}"
            y = "${from.y}"
            this.width = "${width}"
            this.height = "${height}"
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
