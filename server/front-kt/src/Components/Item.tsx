import {Badge, Indicator, Tooltip} from "@mantine/core";

export function Item({name, count}: { name: string, count: number }) {
    const c = count > 1 ? count : null
    // if (itemMap[name]) {
    //     if (!c) {
    //         return <Tooltip label={name} withArrow>
    //             <img src={itemMap[name]} alt={name}/>
    //         </Tooltip>
    //     } else {
    //         return <Tooltip label={name} withArrow>
    //             <Indicator inline label={c} size={12} color="gray" offset={4} radius="md" position="middle-end">
    //                 <img src={itemMap[name]} alt={name}/>
    //             </Indicator>
    //         </Tooltip>
    //     }
    // }
    if (count === 1) {
        return <Badge color="#F0185C" variant="outline" size="sm">{name}</Badge>
    } else {
        return <Badge color="#F0185C" leftSection={count} variant="outline" size="sm">{name}</Badge>
    }
}