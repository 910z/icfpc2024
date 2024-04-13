package com.icfpc.utils

import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper

object Json {
    val mapper = ObjectMapper()

    inline fun <reified T> parse(string: String): T = mapper.readValue(string, T::class.java)
    inline fun <reified T> parse(tree: JsonNode): T = mapper.treeToValue(tree, T::class.java)
    fun <T> toString(obj: T): String = mapper.writeValueAsString(obj)
    fun <T> toObject(obj: T): JsonNode = mapper.valueToTree(obj)
}