//package com.icfpc.problem.model
//
//
//import com.fasterxml.jackson.annotation.JsonProperty
//import com.fasterxml.jackson.databind.ObjectMapper
//import com.icfpc.db.repository.ContentRepository
//import com.icfpc.utils.Json
//
//data class Task(
//    @JsonProperty("room_width")
//    val room_width: Double,
//    @JsonProperty("room_height")
//    val room_height: Double,
//    @JsonProperty("stage_width")
//    val stage_width: Double,
//    @JsonProperty("stage_height")
//    val stage_height: Double,
//    @JsonProperty("stage_bottom_left")
//    val stage_bottom_left: List<Double>,
//    @JsonProperty("musicians")
//    val musicians: List<Int>,
//    @JsonProperty("attendees")
//    val attendees: List<Attendee>,
//    @JsonProperty("pillars")
//    val pillars: List<Pillars>
//) {
//    companion object {
//        fun parse(json: String): Task {
//            return ObjectMapper().readValue(json, Task::class.java)
//        }
//    }
//}
//
//data class Attendee(
//    @JsonProperty("x")
//    val x: Double,
//    @JsonProperty("y")
//    val y: Double,
//    @JsonProperty("tastes")
//    val tastes: List<Double>
//)
//
//data class Pillars(
//    @JsonProperty("center")
//    val center: List<Double>,
//    @JsonProperty("radius")
//    val radius: Double
//)
//
//fun Problem.getContent(contentRepository: ContentRepository) =
//    contentRepository.getReferenceById(contentId).let { Json.parse<Task>(it.content) }
