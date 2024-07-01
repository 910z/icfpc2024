package com.icfpc.utils

import com.icfpc.db.model.History
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Component
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse

@Component
class QueryUtils(
    @Value("\${TOKEN}")
    val token: String
) {
    fun getHistory(): List<History> {
        val request = HttpRequest.newBuilder()
            .uri(URI.create("https://boundvariable.space/team/history?page=1"))
            .header("Accept", "*/*")
            .header("Authorization", "Bearer $token")
            .method("GET", HttpRequest.BodyPublishers.noBody())
            .build()
        val response = HttpClient.newHttpClient().send(request, HttpResponse.BodyHandlers.ofString())
        return Json.parseArray(response.body())
    }

    fun getRequest(uuid: String): String {
        val request = HttpRequest.newBuilder()
            .uri(URI.create("https://boundvariable.space/team/history/$uuid/request"))
            .header("Accept", "*/*")
            .header("Authorization", "Bearer $token")
            .method("GET", HttpRequest.BodyPublishers.noBody())
            .build()
        val response = HttpClient.newHttpClient().send(request, HttpResponse.BodyHandlers.ofString())
        return response.body()
    }

    fun getResponse(uuid: String): String {
        val request = HttpRequest.newBuilder()
            .uri(URI.create("https://boundvariable.space/team/history/$uuid/response"))
            .header("Accept", "*/*")
            .header("Authorization", "Bearer $token")
            .method("GET", HttpRequest.BodyPublishers.noBody())
            .build()
        val response = HttpClient.newHttpClient().send(request, HttpResponse.BodyHandlers.ofString())
        return response.body()
    }
}