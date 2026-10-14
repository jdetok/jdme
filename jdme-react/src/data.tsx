import type { hdrFtrItem, linkDescProps } from "./types";

export const baseUrl = 'https://jdeko.me'

export const headerItems: hdrFtrItem[] = [
    { label: 'jdeko.me', link: baseUrl },
    { label: 'about', link: `${baseUrl}/about` },
    { label: 'source code', link: 'https://github.com/jdetok/jdme' },
    { label: 'supporting tech', link: `${baseUrl}/tech` },
    { label: 'professional resume', link: `${baseUrl}/resume/cv` },
    { label: 'linkedin', link: 'https://www.linkedin.com/in/justin-dekock-257879185' },
];


export const footerItems: hdrFtrItem[] = [
    { label: 'jdeko.me', link: baseUrl },
    { label: 'contact me', link: 'mailto:jdekock17@gmail.com' },
    { label: 'created & maintained by Justin DeKock' },
    { label: 'source code', link: 'https://github.com/jdetok/jdme' },
    { label: 'scroll to top', onClick: () => window.scrollTo({ top: 0, behavior: 'smooth' })},
];

export const prjLinks: linkDescProps[] = [
    {
        label: 'stl transit map | master research project',
        link: 'https://jdeko.me/mrp/',
        desc: 'interactive map webpage built using the ArcGIS Maps SDK for JavaScript from ESRI and self hosted data via Postgres/PostGIS',
    },
    {
        label: 'jdeko.me/bball | nba/wnba stats', 
        link: 'https://jdeko.me/bball/',
        desc: 'browse nba/wnba stats from my nba postgres database using the jdeko.me/bball api written in Go',
    },
    {
        label: 'nba/wnba stats database/ETL', 
        link: 'https://github.com/jdetok/gopuml/',
        desc: 'go cli program & postgres databse designed to fetch & store stats from nba.com',
    },
    {
        label: 'resume/cv', 
        link: 'https://jdeko.me/resume/cv',
        desc: 'my professional resume built with Hugo - a static site generator writtin in Go',
    },
    {
        label: 'slu fleur de lis led circuit', 
        link: 'https://jdeko.me/electro/slu_leds.html',
        desc: 'microcontroller powered led display modeled after Saint Louis University Fleur de Lis logo',
    },
    {
        label: 'brontosaurus led circuit', 
        link: 'https://jdeko.me/electro/bronto.html',
        desc: 'led microcontroller display in the shape of a brontosaurus',
    },
];