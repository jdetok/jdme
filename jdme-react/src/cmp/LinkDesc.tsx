import type { linkDescProps } from "../types";

export default function LinkDesc({ label, link, desc }: linkDescProps) {
    return (
        <div>
            <a target='_blank' href={link}>{label}</a>
            {desc && <p>{desc}</p>}
        </div>
    )
}