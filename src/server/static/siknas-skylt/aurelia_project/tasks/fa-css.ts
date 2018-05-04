import * as gulp from 'gulp';
import * as changedInPlace from 'gulp-changed-in-place';
import * as project from '../aurelia.json';
import {build} from 'aurelia-cli';

export default function faCss() {
    console.log(`PAATH: ${project.paths.fa}/css/*.min.css`)
    return gulp.src(`${project.paths.fa}/css/*.min.css`)
        .pipe(changedInPlace({firstPass:true}))
        // this ensures that our 'require-from' path  
        //  will be simply './font-awesome.min.css'
        .pipe(gulp.dest(project.paths.faCssOutput));
};