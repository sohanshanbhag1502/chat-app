package main

import "github.com/TwiN/go-color"

func main() {
    // Special
    println(color.InBold("This is bold"))
    println(color.InUnderline("This is underlined"))
    // Text colors
    println(color.InBlack("This is in black"))
    println(color.InRed("This is in red"))
    println(color.InGreen("This is in green"))
    println(color.InYellow("This is in yellow"))
    println(color.InBlue("This is in blue"))
    println(color.InPurple("This is in purple"))
    println(color.InCyan("This is in cyan"))
    println(color.InGray("This is in gray"))
    println(color.InWhite("This is in white"))
    // Background colors
    println(color.OverBlack("This is over a black background"))
    println(color.OverRed("This is over a red background"))
    println(color.OverGreen("This is over a green background"))
    println(color.OverYellow("This is over a yellow background"))
    println(color.OverBlue("This is over a blue background"))
    println(color.OverPurple("This is over a purple background"))
    println(color.OverCyan("This is over a cyan background"))
    println(color.OverGray("This is over a gray background"))
    println(color.OverWhite("This is over a white background"))
}