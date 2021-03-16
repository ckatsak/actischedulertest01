// Package hardcoded contains an implementation of acti.InterferenceModel,
// where the slowdowns among all applications are known and hardcoded.
package hardcoded

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

///////////////////////////////////////////////////////////////////////////////
//
// HardcodedSlowDowns
//
///////////////////////////////////////////////////////////////////////////////

// HardcodedSlowDowns is an implementation of acti.InterferenceModel, where the
// slowdowns among all applications are known and hardcoded.
type HardcodedSlowDowns struct {
	actiLabelKey string
}

// New returns a new HardcodedSlowDowns with the given label key (the one that
// is used by ActiPlugin to track its applications).
func New(actiLabelKey string) *HardcodedSlowDowns {
	return &HardcodedSlowDowns{
		actiLabelKey: actiLabelKey,
	}
}

// Attack implements acti.InterferenceModel; see the documentation there for
// more information.
func (m *HardcodedSlowDowns) Attack(attacker, occupant *corev1.Pod) (float64, error) {
	occPodCategory, _ := parseAppCategory(occupant.Labels[m.actiLabelKey])
	//  occupant   ^^^   Pod's label's value must have been
	// validated back when it was first scheduled on the Node
	newPodCategory, err := parseAppCategory(attacker.Labels[m.actiLabelKey])
	if err != nil {
		return -1, err
	}
	return newPodCategory.attack(occPodCategory), nil
}

const toInt64Multiplier = 100.

// ToInt64Multiplier implements acti.InterferenceModel; see the documentation
// there for more information.
func (_ *HardcodedSlowDowns) ToInt64Multiplier() float64 {
	return toInt64Multiplier
}

///////////////////////////////////////////////////////////////////////////////
//
// appCategory
//
///////////////////////////////////////////////////////////////////////////////

// appCategory is an enumeration of known application categories.
type appCategory int64

const (
	catA appCategory = iota
	catB
	catC
	catD
	catE
)

// String returns the string representation of the (known) appCategory.
func (ac appCategory) String() string {
	switch ac {
	case catA:
		return "catA"
	case catB:
		return "catB"
	case catC:
		return "catC"
	case catD:
		return "catD"
	case catE:
		return "catE"
	default:
		return "UNKNOWN"
	}
}

// parseAppCategory parses a (known) appCategory from a string.
func parseAppCategory(category string) (appCategory, error) {
	switch category {
	case "catA":
		return catA, nil
	case "catB":
		return catB, nil
	case "catC":
		return catC, nil
	case "catD":
		return catD, nil
	case "catE":
		return catE, nil
	default:
		return -1, fmt.Errorf("unknown application category: '%s'", category)
	}
}

// attack returns the slowdown incurred on the given occupant when the
// appCategory is scheduled along with it.
func (ac appCategory) attack(occupant appCategory) float64 {
	return slowDowns[ac][occupant]
}

///////////////////////////////////////////////////////////////////////////////
//
// slowDownMatrix
//
///////////////////////////////////////////////////////////////////////////////

// slowDownMatrix is a type alias for internal use in ActiPlugin.
type slowDownMatrix map[appCategory]map[appCategory]float64

// slowDowns is a hardcoded global map that represents a dense 2D matrix of the
// slowdowns incurred by application colocations. Its format is as follows:
//
//     {
//         A: {
//             A: f64 slowdown of an A when attacked by an A
//             B: f64 slowdown of a B when attacked by an A
//             C: f64 slowdown of a C when attacked by an A
//         },
//         B: {
//             A: f64 slowdown of an A when attacked by a B
//             B: f64 slowdown of a B when attacked by a B
//             C: f64 slowdown of a C when attacked by a B
//         },
//         C: {
//             A: f64 slowdown of an A when attacked by a C
//             B: f64 slowdown of a B when attacked by a C
//             C: f64 slowdown of a C when attacked by a C
//         },
//         . . .
//     }
var slowDowns = slowDownMatrix{
	catA: map[appCategory]float64{
		catA: 1.15,
		catB: 2.43, // slowdown of catB when attacked by catA = 2.43
		catC: 1.22,
		catD: 1.34,
		catE: 1.41,
	},
	catB: map[appCategory]float64{
		catA: 2.19,
		catB: 1.26,
		catC: 1.48,
		catD: 3.10, // slowdown of catD when attacked by catB = 3.10
		catE: 1.37,
	},
	catC: map[appCategory]float64{
		catA: 1.21,
		catB: 3.09,
		catC: 2.12, // slowdown of catC when attacked by catC = 2.12
		catD: 1.18,
		catE: 1.53,
	},
	catD: map[appCategory]float64{
		catA: 1.47,
		catB: 1.24,
		catC: 1.76,
		catD: 2.25,
		catE: 1.10, // slowdown of catE when attacked by catD = 1.10
	},
	catE: map[appCategory]float64{
		catA: 3.15, // slowdown of catA when attacked by catE = 3.15
		catB: 1.74,
		catC: 1.16,
		catD: 1.27,
		catE: 2.53,
	},
}
