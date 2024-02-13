<template>
    <div class="container">
        <div class="mb-5">
            <div class="field has-addons">
                <p class="control">
                    <a :class="['button', 'is-static', 'is-large']">
                        <!-- attempts/allowed_attempts -->
                        3/3
                    </a>
                </p>
                <p class="control is-expanded has-icons-left">
                    <input :class="['input', 'is-large']" ref="input" :type="getInputType" placeholder="Bib Number" />
                    <span class="icon is-large is-left">
                        <icon :icon="['fas', 'fa-user']"></icon>
                    </span>
                </p>
                <p class="control" v-if="false">
                    <a :class="['button', 'is-static', 'is-large']">
                        <!-- best attempt -->
                        <!-- {{ time }} -->
                    </a>
                </p>
            </div>
            <label class="checkbox">
                <input type="checkbox" v-model="letterInput">
                Allow Letters
            </label>
        </div>
        <div class="mb-5">
            <div class="field has-addons">
                <div class="control">
                    <div class="dropdown is-hoverable">
                        <div class="dropdown-trigger">
                            <button class="button is-large is-info" aria-haspopup="true"
                                aria-controls="distance-measuring-menu">
                                <span>Distance Measuring</span>
                                <span class="icon is-small">
                                    <icon :icon="['fas', 'fa-angle-down']" aria-hidden="true"></icon>
                                </span>
                            </button>
                        </div>
                        <div class="dropdown-menu" id="distance-measuring-menu" role="menu">
                            <div class="dropdown-content">
                                <a href="#" class="dropdown-item">
                                    Feet & Inches
                                </a>
                                <a href="#" class="dropdown-item">
                                    Meters & Centimeters
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="control is-expanded has-icons-left">
                    <input class="input is-large" placeholder="Feet" type="number">
                    <span class="icon is-large is-left">
                        <icon :icon="['fas', 'fa-ruler']"></icon>
                    </span>
                </div>
                <div class="control is-expanded has-icons-left">
                    <input class="input is-large" placeholder="Inches" type="number">
                    <span class="icon is-large is-left">
                        <icon :icon="['fas', 'fa-ruler']"></icon>
                    </span>
                </div>
                <div class="control">
                    <a class="button is-large is-danger">Record!</a>
                </div>
            </div>
        </div>
        <div class="level">
            <div class="level-item has-text-centered">
                <div>
                    <p class="heading">
                        <icon :icon="['fas', 'fa-crown']"></icon> Attempt #2
                    </p>
                    <p class="title">5'5"</p>
                </div>
            </div>
            <div class="level-item has-text-centered">
                <div>
                    <p class="heading">Attempt #1</p>
                    <p class="title">5'4"</p>
                </div>
            </div>
            <div class="level-item has-text-centered">
                <div>
                    <p class="heading">Attempt #3</p>
                    <p class="title">5'2"</p>
                </div>
            </div>
        </div>
        <!-- <p :class="['help']" v-if="result.show">{{ result.msg }}</p> -->
    </div>
</template>

<script>
import { setAuthHeader } from '../api/auth';

export default {
    mounted: function () {
    },
    data: function () {
        return {
            letterInput: false,
            attempts: [
                {
                    attempt: 3,
                    distance: 120,
                },
                {
                    attempt: 1,
                    distance: 123,
                },
                {
                    attempt: 2,
                    distance: 119,
                }
            ]
        }
    },
    computed: {
        getInputType() {
            return this.letterInput ? "text" : "number"
        }
    },
    methods: {
        recordAttempt: function () {

        },
    }
};
</script>
