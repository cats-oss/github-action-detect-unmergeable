'use strict';

const assert = require('assert');

const { hasNeedRebaseLabel, LABEL_STATUS_NEED_REBASE } = require('./labels');

async function getDefaultBranchName(octokit, owner, repo) {
    const result = await octokit.repos.get({owner, repo});
    const body = result.data;
    assert.strictEqual(!!body, true);

    const branchName = body.default_branch;
    return branchName;
}

async function getOpenPullRequestAll(octokit, owner, repo) {
    // FIXME: this code could not get all opened pull requests if there are over 100.
    const result = await octokit.pulls.list({
        owner,
        repo,
        state: 'open',
        // eslint-disable-next-line camelcase
        per_page: '100',
    });
    const body = result.data;
    assert.strictEqual(!!body, true);

    return body;
}

function getPullRequestId(pullReqInfo) {
    const id = pullReqInfo.number;
    assertTypeof(id, 'number');
    return id;
}

const SECONDS_WE_ASSUME_GITHUB_WOULD_COMPLETE_CHECK_IF_MERGEBLE = 1000 * 10;

async function checkAndMarkIfPullRequestUnmergeable(octokit, owner, repo, oldPullReqInfo, compareUrl) {
    const number = getPullRequestId(oldPullReqInfo);
    if (hasNeedRebaseLabel(oldPullReqInfo.labels)) {
        console.log(`#${number} has been labeled as ${LABEL_STATUS_NEED_REBASE}`);
        return;
    }

    // When we get all opened pull requests, GitHub have not checked whether ths PR is unmergeble yet.
    // So I think we should retry them.
    let mergeable = await shouldMarkPullRequestNeedRebase(octokit, owner, repo, number);
    if (mergeable === null) {
        // retry
        await sleep(SECONDS_WE_ASSUME_GITHUB_WOULD_COMPLETE_CHECK_IF_MERGEBLE);
        mergeable = await shouldMarkPullRequestNeedRebase(octokit, owner, repo, number);
    }

    console.log(`The mergeable of #${number} is ${mergeable}`);
    if (mergeable === null) {
        // give up
        return;
    }

    if (mergeable) {
        return;
    }

    const addLabel = octokit.issues.addLabels({
        owner,
        repo,
        number,
        labels: [LABEL_STATUS_NEED_REBASE],
    });

    const body = `:umbrella: The latest upstream change (presumably [these](${compareUrl})) made this pull request unmergeable. Please resolve the merge conflicts.`;
    const addComment = octokit.issues.createComment({
        owner,
        repo,
        number,
        body,
    });

    await Promise.all([addLabel, addComment]);
}

async function shouldMarkPullRequestNeedRebase(octokit, owner, repo, number) {
    const newPRInfoResponse = await octokit.pulls.get({
        owner,
        repo,
        number,
    });
    const newPRInfo = newPRInfoResponse.data;
    assert.strictEqual(!!newPRInfo, true);

    // Check again to confirm the other instance of this action's behavior.
    if (hasNeedRebaseLabel(newPRInfo.labels)) {
        return false;
    }

    const mergeable = newPRInfo.mergeable;
    if (mergeable === null) {
        return null;
    }

    return mergeable;
}

function sleep(millisec) {
    const p = new Promise((resolve) => {
        setTimeout(resolve, millisec);
    });
    return p;
}

function assertTypeof(val, typename, message) {
    assert.strictEqual(typeof val, typename, message);
}

module.exports = Object.freeze({
    getDefaultBranchName,
    getOpenPullRequestAll,
    checkAndMarkIfPullRequestUnmergeable,
});
