<script>
import {h} from 'vue'
import CondItem from './CondItem.vue'

let renderNode = (node, readonly) => {
    if (node.type == "and" || node.type == "or") {
        return h(CondItem, {node: node, readonly: readonly}, () => {
            let ret = []
            node.subs.forEach(value => {
                ret.push(renderNode(value, readonly))
            });
            return ret;
        })
    }
    return h(CondItem, {node: node, readonly: readonly}, null)
}

export default {
    props: ['tree', 'readonly'],
    setup(props) {
        return () => renderNode(props.tree, props.readonly)
    }
}
</script>