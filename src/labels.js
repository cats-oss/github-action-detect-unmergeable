'use strict';

const STATUS_LABEL_PREFIX = 'S-';
const LABEL_STATUS_NEED_REBASE = `${STATUS_LABEL_PREFIX}needs-rebase`;

/**
 *  @param {!Array<{ name: string; }>} input
 *  @returns    {boolean}
 */
function hasNeedRebaseLabel(input) {
    return input
        .map((obj) => obj.name)
        .some((label) => {
            return label === LABEL_STATUS_NEED_REBASE;
        });
}

module.exports = Object.freeze({
    hasNeedRebaseLabel,
    LABEL_STATUS_NEED_REBASE,
});
