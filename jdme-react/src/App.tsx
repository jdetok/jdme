import Footer from "./cmp/Footer";
import Header from "./cmp/Header";
import LinkDesc from "./cmp/LinkDesc";
import * as d from './data';
export default function App() {
    return (
        <>
            <Header cssClass='hdr' items={d.headerItems} />
            <main>
                <h1>jdeko.me</h1>
                <h2 className='h2-1'>Justin DeKock's personal website/portfolio - hosted locally on a Raspberry Pi 5</h2>
                <h2 className='h2-2'>programming/system design projects</h2>
                <div className='projects'>
                    {d.prjLinks.map((item, i) => (
                        <LinkDesc key={i} { ...item } />
                    ))}
                </div>
                <p>/-|-\</p>
                <h2 className='h2-2'>listen to my music</h2>

            </main>
            <Footer cssClass='ftr' items={d.footerItems} />
        </>
    )
};
