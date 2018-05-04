import * as gulp from 'gulp';
import * as project from '../aurelia.json';

export default function faFonts() {
    return gulp.src(`${project.paths.fa}/fonts/*`)
        .pipe(gulp.dest(project.paths.faFontsOutput));
}