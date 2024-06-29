package com.icfpc.lang

import java.math.BigInteger
object Util {
    val line = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`|~ \n"

    fun parse(command: String): Any? {
        if (!command.startsWith("S")) return null
        var num = BigInteger.ZERO
        for (i in command.drop(1)) {
            num = num * 94.toBigInteger() + (i.code - 33).toBigInteger()
        }
        var s = ""
        while (num > BigInteger.ZERO) {
            val a = num % line.length.toBigInteger()
            s = line[a.toInt()] + s
            num /= line.length.toBigInteger()
        }
        return s
    }
}