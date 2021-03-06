package artemis_tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMovie_AddNameNoUnderscores(t *testing.T) {
	expected := "a_river_runs_through_it_brad_pitt.mp4"

	n := AddName("a_river_runs_through_itbradpitt.mp4", "")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameUnknownWithBrackets(t *testing.T) {
	expected := "scene_480p_1_brad_pitt.mp4"

	n := AddName("scene_480p (1).mp4", "")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameNoPrecedingUnderscore(t *testing.T) {
	expected := "a_river_runs_through_it_brad_pitt.mp4"

	n := AddName("a_river_runs_through_it_bradpitt.mp4", "")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameAtStartNoPrecedingUnderscore(t *testing.T) {
	expected := "brad_pitt_a_river_runs_through_it.mp4"

	n := AddName("bradpitta_river_runs_through_it.mp4", "")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameAtStartMiddleUnderscore(t *testing.T) {
	expected := "brad_pitt_a_river_runs_through_it.mp4"

	n := AddName("brad_pitta_river_runs_through_it.mp4", "")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameNoUnderscoresWithMiddleName(t *testing.T) {
	expected := "a_river_runs_through_it_brad_tiberius_pitt.mp4"

	n := AddName("a_river_runs_through_itbradtiberiuspitt.mp4", "tiberius")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameNoPrecedingUnderscoreWithMiddleName(t *testing.T) {
	expected := "a_river_runs_through_it_brad_tiberius_pitt.mp4"

	n := AddName("a_river_runs_through_it_bradtiberius_pitt.mp4", "tiberius")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameAtStartNoPrecedingUnderscoreWithMiddleName(t *testing.T) {
	expected := "brad_tiberius_pitt_a_river_runs_through_it.mp4"

	n := AddName("bradtiberiuspitta_river_runs_through_it.mp4", "tiberius")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameAtStartMiddleUnderscoreWithMiddleName(t *testing.T) {
	expected := "brad_tiberius_pitt_a_river_runs_through_it.mp4"

	n := AddName("bradtiberius_pitta_river_runs_through_it.mp4", "tiberius")

	assert.Equal(t, expected, n, "Movie names should match")
}

func TestMovie_AddNameNoUnderscoresMultipleActors(t *testing.T) {
	expected := "a_river_runs_through_it_brad_pitt_robert_redford.mp4"
	nms := []string{"brad_pitt", "robert_redford"}
	m := CreateMovieWithActors("a_river_runs_through_it_brad_pittrobertredford.mp4", &nms)
	m.AddActorNames()

	assert.Equal(t, expected, m.NewName, "Movie names should match")
}

func TestMovie_AddNameNoUnderscoresMultipleActorsWithMiddleName(t *testing.T) {
	expected := "a_river_runs_through_it_brad_tiberius_pitt_robert_redford.mp4"
	nms := []string{"brad_tiberius_pitt", "robert_redford"}
	m := CreateMovieWithActors("a_river_runs_through_it_brad_pittrobertredford.mp4", &nms)
	m.AddActorNames()

	assert.Equal(t, expected, m.NewName, "Movie names should match")
}
