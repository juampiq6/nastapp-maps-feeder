package main

// earthRadius = 6378137
// 111139 m = 1 degree

// 0.008997 rad degrees roughly equals 1000m
const differentialRad = 0.0397

// 1000m would be the increment every time

// Diagonal line with 2 points represents square for Mendoza
// this is the top left point
var PointAMendoza = LatLong{Lat: -32.052662, Long: -70.25}

// this is the bottom right point
var PointBMendoza = LatLong{Lat: -37.601944, Long: -65.081140}
